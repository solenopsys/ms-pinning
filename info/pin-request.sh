curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: 3e234212423423423423423423432" \
     -d '{
           "cids": [
             "QmcTr47jWdVKzc9USNXguJKbigSbrcKXj7wi93xHAaeytG"
           ],
           "rep_min": 2,
           "rep_max": 3
         }' \
     http://localhost:8085/pin