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

import alabaster
import sphinxswagger


project = 'sphinx-swagger'
copyright = '2016, Dave Shawley'
release = '.'.join(str(v) for v in sphinxswagger.version_info[:2])
version = sphinxswagger.__version__
needs_sphinx = '1.0'
extensions = [
    'sphinx.ext.intersphinx',
]

master_doc = 'index'
html_theme = 'alabaster'
html_theme_path = [alabaster.get_path()]
html_sidebars = {
    '**': ['about.html',
           'navigation.html'],
}
html_theme_options = {
    'description': 'Generate swagger definitions',
    'github_user': 'dave-shawley',
    'github_repo': 'sphinx-swagger',
}
intersphinx_mapping = {
    'python': ('https://docs.python.org/3', None),
}
