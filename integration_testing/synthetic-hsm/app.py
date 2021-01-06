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
