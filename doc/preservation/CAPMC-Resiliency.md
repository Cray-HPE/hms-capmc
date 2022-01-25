1.  [CASMHMS](index.html)
2.  [CASMHMS Home](CASMHMS-Home_119901124.html)
3.  [Design Documents](Design-Documents_127906417.html)
4.  [CAPMC](CAPMC_227479290.html)

# <span id="title-text"> CASMHMS : CAPMC Resiliency </span>

Created by <span class="author"> David Laine</span>, last modified on
Apr 22, 2020

This document outlines the plan in the Shasta 1.3 - 1.4 timeframe for
insuring service resiliency for the capmc functionality.

If dependent services (state manager, vault, telemetry database) become
unavailable capmc will not be able to function.  This document only
addresses the capmc service.

## Shasta v1.3

#### Definition of resiliency for this release

For this release the capmc service will always be available to call. 
There won't be times that the service is down, restarting, or hung. 
Other services will always be able to make calls to it, but the calls
are not insured to complete correctly.  

#### Implications

If the server an instance of capmc is running on becomes unavailable,
requests will be routed to other replicas of the pod and continue
working while that replica gets restored.

When a call is made to capmc, it is handled by the pod that the call is
directed to.  That pod alone knows about the call and handles it
completely with no communication or shared storage between capmc pods. 
If the pod crashes or stops before that action is complete, it is not
picked up by anything else and all knowledge of state is lost when the
pod is restarted.

Many single calls to capmc result in actions applied to multiple nodes. 
A single power on/off/reset command can be distributed to tens,
hundreds, or even thousands of individual nodes.  If the pod doing the
action stops part way through there is no way for anything else to
either pick up and complete the work, or even know what portion was
completed and what still remains undone.

The end result is that users may need to be aware of incomplete
operations and have the awareness of needing to restart or rerun an
operation that does not complete cleanly.  The caller will know this
happens due to a broken connection with the capmc service rather than a
valid return from the call.

Upon investigation, if a pod dies mid-command (is killed and restarted
by kubernetes) the command is automatically re-executed on another
replica of the pod.  There is no history kept, so the command is started
over from the beginning.

#### Plan to insure this definition of resiliency

To achieve this level of resiliency all that is required is to scale up
the number of replicas to 3.  There is no state contained in capmc so
having multiple replicas should not pose an issue with any of the
current workflows of capmc.  The locking in state manager will prevent
multiple instances from trying to operate on the same endpoint at the
same time - given the current parallelism already inside capmc this is
an issue even with a single instance of the pod in operation.

#### Required work

1.  Fine tune the state manager locking to shorten lock times. ( <span
    class="jira-issue resolved" jira-key="CASMHMS-3282">
    <a href="https://connect.us.cray.com/jira/browse/CASMHMS-3282?src=confmacro" class="jira-issue-key"><img src="https://connect.us.cray.com/jira/secure/viewavatar?size=xsmall&amp;avatarId=13311&amp;avatarType=issuetype" class="icon" />CASMHMS-3282</a>
    - <span class="summary">CAPMC: Fine tune component locking via state
    manager</span> <span
    class="aui-lozenge aui-lozenge-subtle aui-lozenge-success jira-macro-single-issue-export-pdf">Closed</span>
    </span> )
    1.  This looks pretty invasive - better to just leave as-is in capmc
        and build this into PCS - start PCS work that much sooner?
2.  Increase replica count and testing. ( <span
    class="jira-issue resolved" jira-key="CASMHMS-333">
    <a href="https://connect.us.cray.com/jira/browse/CASMHMS-333?src=confmacro" class="jira-issue-key"><img src="https://connect.us.cray.com/jira/secure/viewavatar?size=xsmall&amp;avatarId=13315&amp;avatarType=issuetype" class="icon" />CASMHMS-333</a>
    - <span class="summary">cray-capmc resiliency (HA)</span> <span
    class="aui-lozenge aui-lozenge-subtle aui-lozenge-success jira-macro-single-issue-export-pdf">Done</span>
    </span> )
    1.  A change is being made to release active locks when k8s kills a
        pod, that way when the command is re-issued to another pod it
        will not immediately die for lack of ability to obtain the lock.

#### Further investigation

We need to look into the calls that may be most susceptible to
mid-flight incompletions and the services that may be calling these and
how they would react.

1.  Given a 'reset' of multiple endpoints is there a way to find which
    were completed and where to pick up?
    1.  ANSWER: No - isn't possible to figure out which were reset and
        which were not.  Given the current system the entire call would
        need to be redone until it succeeds.
2.  What other services interact with capmc directly and how do they
    recover from an error in a call?
    1.  hms-mountain-discovery - query power status and turn on nodes
        that are off (single attempt, no retries)
    2.  boa/bss - one attempt and no retry
    3.  ??? others?
3.  If capmc dies in the middle of a command can the same command (with
    identical arguments) be called again without error?
    1.  example: command 'all on' given and dies 1/2 way through, will
        'on' command called again result in error since now some are
        already on?
    2.  ANSWER: hardware dependent - some hardware throws an error if
        already in requested state, some does nothing.

## Shasta v1.4

#### Definition of resiliency for this release

For this release capmc will be call resilient where a call once made is
insured to be carried out to completion and without duplicate commands
to endpoints no matter what happens to the individual pods doing the
work.  This is the definition of resiliency we need to achieve for power
operations. 

#### <span style="color: rgb(0,0,0);">Plan to insure this definition of resiliency</span>

There are several parts that need to be in place for this level of
resiliency to be achieved.  These may be built into capmc, created as
separate services that capmc uses, or built into a new service where
capmc just becomes the back-compatible API into the new service.

1\) We need a truly RESTful API.  As the number of nodes increases it
will become more difficult to complete an operation within the standard
time window for an http call.  Many times the hardware interaction
itself is too slow for a 30 or 60 second timeout for http interactions. 
We need to provide an API where the user is expecting to initiate an
operation and get a key back where they are expected to check back in
for the progress and completion of the operation.  That allows for long
times to complete the operation, intermediate status updates, a
mid-operation cancel, and large packages of detailed return
information.  The current capmc API is inadequate for these situations.

2\) We need a concept of a shared worker pool that is itself resilient
to failures. 

-   The initial command into capmc will generate the set of commands to
    be executed on the actual nodes / endpoints of the system.  
-   These commands will need to be stored as part of a single job where
    multiple workers will need to pull individual commands off of the
    queue, mark 'in progress' states, and eventual completion
    information.  
-   There will need to be an access point where the user can call to
    find the state of the job and query eventual results.  
-   There may be timeout conditions built in where an individual command
    is retried if it has been determined the worker failed
    mid-command.  
-   This requirement should be broken out into a separate design
    document when the time comes to create it.  
-   NOTE: many services may require this kind of functionality -
    building on what is learned from current FAS work, it would be good
    if this could be built into a common tool that multiple services
    could use.

3\) We need something to monitor the dependent services more closely and
figure out the implications for how these commands would react when
other services are down.

  

  

## Comments:

<table data-border="0" width="100%">
<colgroup>
<col style="width: 100%" />
</colgroup>
<tbody>
<tr class="odd">
<td><span id="comment-172197136"></span>
<p>Should we be calling out PCS for 1.4 </p>
<p>When we state CAPMC is not restful we should state the reason why we didn't do that in the first the place and the answer can be that we cloned the functionality of XC CAPMC since that was what was decide long ago.     So we are where we are but should call that out.</p>
<div class="smallfont" data-align="left" style="color: #666666; width: 98%; margin-bottom: 10px;">
<img src="images/icons/contenttypes/comment_16.png" width="16" height="16" /> Posted by rfrost at Apr 07, 2020 15:49
</div></td>
</tr>
<tr class="even">
<td style="border-top: 1px dashed #666666"><span id="comment-172197145"></span>
<p>Thats a good point, I think we should probably call that out a bit more.  <br />
<br />
Also worth pointing out is that we need to do a requirements validation with our stakeholders... CAPMC in shasta tried to force the cascade paradigm on Shasta.  While that has 'worked' I dont think it is necessarily appropriate for what Shasta needs.  In other words, while we understand some of the limitations, or tradeoffs of CAPMC, we need to make sure we understand what PCS really needs to do before we make a 'RESTful' CAPMC.  </p>
<div class="smallfont" data-align="left" style="color: #666666; width: 98%; margin-bottom: 10px;">
<img src="images/icons/contenttypes/comment_16.png" width="16" height="16" /> Posted by anieuwsma at Apr 07, 2020 15:57
</div></td>
</tr>
</tbody>
</table>

Document generated by Confluence on Jan 14, 2022 07:17

[Atlassian](http://www.atlassian.com/)
