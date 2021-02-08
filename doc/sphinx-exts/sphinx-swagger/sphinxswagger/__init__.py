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

version_info = (0, 0, 4)
__version__ = '.'.join(str(v) for v in version_info)


def setup(app):
    """
    Called by Sphinx to initialize the extension.

    :param sphinx.application.Sphinx app: sphinx instance that is running
    :return: a :class:`dict` of extension metadata -- ``version`` is the
        only required key
    :rtype: dict

    """
    from . import builder, writer

    app.add_builder(builder.SwaggerBuilder)
    app.add_config_value('swagger_file', 'swagger.json', True)
    app.add_config_value('swagger_license', {'name': 'Proprietary'}, True)
    app.add_config_value('swagger_description', '', True)
    app.connect('build-finished', writer.write_swagger_file)

    return {'version': __version__}
