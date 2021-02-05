#  MIT License
#
#  (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
#
#  Permission is hereby granted, free of charge, to any person obtaining a
#  copy of this software and associated documentation files (the "Software"),
#  to deal in the Software without restriction, including without limitation
#  the rights to use, copy, modify, merge, publish, distribute, sublicense,
#  and/or sell copies of the Software, and to permit persons to whom the
#  Software is furnished to do so, subject to the following conditions:
#
#  The above copyright notice and this permission notice shall be included
#  in all copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
#  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
#  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
#  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
#  OTHER DEALINGS IN THE SOFTWARE.

from flask import Flask, render_template, request, jsonify, make_response
import json

app = Flask(__name__)

@app.route('/hsm/v1/service/ready', methods=['GET'])
def service_ready():
    return jsonify(code=0, message="Ready")

@app.route('/hsm/v1/State/Components', methods=['POST','GET'])
def state_components():
    with open('files/state_components.json', 'r') as f:
        components = json.load(f)

    nidList = request.args.getlist('nid')
    nids = []

    #convert from string to int
    for n in nidList :
        try:
            nid = int(n)
            nids.append(nid)
        except ValueError:
            continue

    #if there are no nids in the query then return all of them from the file
    if len(nids) == 0 :
        return jsonify(components)

    #else match the nids and return
    componentsByNid = {}
    for c in components["Components"]:
        componentsByNid[c["NID"]] = c

    nidComponents = []
    for nid in nids:
        if nid in componentsByNid.keys():
            nidComponents.append(componentsByNid[nid])

    finalComponents = {}
    finalComponents["Components"] = nidComponents

    return jsonify(finalComponents)


if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=27779)
