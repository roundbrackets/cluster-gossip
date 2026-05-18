# cluster-gossip

k8s gossip listens for k8s events and publishes summaries to a database beackend for a curious person to review.

What it does:
* Listens to events
* Sends it to an LLM which creates a summary with a type and frequency
* The summary is written to a database
* The summaries are available via a webserver
* The summaries are sent to telegram

What it does not:
* It does not act on errors or attempt to perform any type of troubleshooting
* It does not attempt to put the events in a broader context of all the events previously logged

What it is:
- a deployment on the k8s cluster itself which talks to an local model of a server outside of the cluster. 

What it is not:

Cluster Gossiper is not splunk, grafana, or any number of commercial and open source monotpring project. Those products provide much more comprehensive fuctionality. I am not trying to replace or compete with them, but rather create a lightweight intospection into my own respberry pi k3s cluster.


## References
* https://etcd.io/docs/v3.2/learning/api/
* https://docs.k3s.io/
