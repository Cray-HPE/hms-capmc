@Library('dst-shared@master') _

dockerBuildPipeline {
        githubPushRepo = "Cray-HPE/hms-capmc"
        repository = "cray"
        imagePrefix = "cray"
        app = "capmc"
        name = "hms-capmc"
        description = "Cray CAPMC service"
        dockerfile = "Dockerfile"
        slackNotification = ["", "", false, false, true, true]
        product = "csm"
}
