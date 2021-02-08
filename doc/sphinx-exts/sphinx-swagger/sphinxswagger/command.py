# MIT License
#
# (C) Copyright [2021] Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

from distutils import cmd, log
import os.path

from sphinx import application


class BuildSwagger(cmd.Command):
    description = 'Build a swagger definition from Sphinx docs'
    user_options = [
        ('config-dir=', 'c', 'configuration directory'),
        ('output-file=', 'o', 'output file name'),
        ('ignore-distinfo', 'u', 'ignore distribution metadata'),
    ]
    boolean_options = ['ignore-distinfo']

    def initialize_options(self):
        self.config_dir = None
        self.output_file = None
        self.ignore_distinfo = False

    def finalize_options(self):
        if self.config_dir is None:
            self.config_dir = 'docs'
        self.ensure_dirname('config_dir')
        if self.config_dir is None:
            self.config_dir = os.curdir
            self.warning('Using {} as configuration directory',
                         self.source_dir)
        self.config_dir = os.path.abspath(self.config_dir)

        if self.output_file is not None:
            self.output_file = os.path.abspath(self.output_file)

    def run(self):
        build_cmd = self.get_finalized_command('build')
        build_dir = os.path.join(os.path.abspath(build_cmd.build_base),
                                 'swagger')
        self.mkpath(build_dir)
        doctree_dir = os.path.join(build_dir, 'doctrees')
        self.mkpath(doctree_dir)

        overrides = {}
        if self.output_file is not None:
            overrides['swagger_file'] = self.output_file

        if not self.ignore_distinfo:
            if self.distribution.get_description():
                overrides['swagger_description'] = \
                    self.distribution.get_description()
            if self.distribution.get_license():
                overrides['swagger_license.name'] = \
                    self.distribution.get_license()
            if self.distribution.get_version():
                overrides['version'] = self.distribution.get_version()

        app = application.Sphinx(
            self.config_dir, self.config_dir, build_dir, doctree_dir,
            'swagger', confoverrides=overrides)
        app.build()

    def warning(self, msg, *args):
        self.announce(msg.format(*args), level=log.WARNING)

    def info(self, msg, *args):
        self.announce(msg.format(*args), level=log.INFO)

    def debug(self, msg, *args):
        self.announce(msg.format(*args), level=log.DEBUG)
