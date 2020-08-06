# Scripts

- `create_service.sh`
  - creates managed instance group and backend service. They are the VMs that will serve gRPC traffic.
- `create_health_check.sh`
  - creates a gRPC health check, and a firewall rule to allow the health check traffic. This allows the backend service to know which backend VMs are healthy and ready to serve traffic.
- `config_traffic_director.sh`
  - creates traffic director configutation, including URL-map and forwarding rule using the URL-map. The traffic director will send those to the gRPC clients to control how traffic will be routed.

- `all.sh`
  - calls scripts above to create all resources
- `clean.sh`
  - deletes all resources created