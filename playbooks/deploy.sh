#!/bin/bash

# Inventory file..
INVENTORY_FILE="inventory.yaml"

# Run the playbooks..
ansible-playbook -i "$INVENTORY_FILE" "./build.yaml"
ansible-playbook -i "$INVENTORY_FILE" "./01 - build_resources.yaml"
ansible-playbook -i "$INVENTORY_FILE" "./02 - get_certificates.yaml"

# Deploy the application..
ansible-playbook -i "$INVENTORY_FILE" "deploy.yaml"
