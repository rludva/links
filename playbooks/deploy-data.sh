#!/bin/bash

# Inventory file..
INVENTORY_FILE="inventory.yaml"

# Deploy the data..
ansible-playbook -i "$INVENTORY_FILE" "deploy-data.yaml" -e "restart_service=true"
