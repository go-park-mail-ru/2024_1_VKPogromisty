#!/bin/bash
sudo -i -u postgres
psql -h localhost -p 5432 -U postgres -d socio
