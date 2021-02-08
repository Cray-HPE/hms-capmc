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

import docutils.io

from sphinx import builders

from . import document, writer


class SwaggerBuilder(builders.Builder):
    name = 'swagger'
    allow_parallel = False

    def init(self):
        """Sub-class hook called from __init__"""
        self.writer = None

    def prepare_writing(self, docnames):
        """Called before :meth:`write_doc`"""
        self.swagger = document.SwaggerDocument()
        self.writer = writer.SwaggerWriter(swagger_document=self.swagger)

    def write_doc(self, docname, doctree):
        """Write a doc to the filesystem."""
        destination = docutils.io.NullOutput()
        self.writer.write(doctree, destination)

    def get_outdated_docs(self):
        """List of docs that we need to write or just a file name."""
        return self.app.config.swagger_file

    def get_target_uri(self, docname, typ=None):
        return ''  # No clue what to return here :/

    def finish(self):
        """Called after write() has completed."""
        pass
