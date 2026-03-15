#!/bin/sh

/app/vikunja/vikunja user create --email 'my-home@default.com' -p 'default' -u 'Home'

exec /app/vikunja/vikunja
