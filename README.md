# Scribe Postman Newman InfluxDB Go (scribe-postmannewman-influxdb-go)

Scribe to read data from InfluxDB and show in Craftsperson

---
## Getting Started
### Requirements
1. Docker or K8s
2. InfluxDB2 [[Docker](https://github.com/Sea-Creeper/ansible-playbooks/tree/main/docker/influxdb) | [K8s]()]
3. Environment variables:
   1. INFLUXDB_ORG
   2. INFLUXDB_TOKEN
4. Run config
   1. Network `db_creeper`
   2. Port 9080:9080
---
### Ansible
Docker Deploy
```
ansible-playbook -i ansible/hosts ansible/docker/deploy.yml
```

Docker Down
```
ansible-playbook -i ansible/hosts ansible/docker/down.yml 
```

---

## Scribe fullname format
`scribe <creeper> <datasource> <language>`

Example,
`scribe postmannewman influxdb go`