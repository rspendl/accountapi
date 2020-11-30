curl -v dockerw0:8080/v1/organisation/accounts -H "Content-Type: application/vnd.api+json" -d'
{
  "data": {
    "type": "accounts",
    "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
    "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    "attributes": {
      "country": "GB"
    }
  }
}'
exit

curl -v dockerw0:8080/v1/organisation/accounts -H "Content-Type: application/vnd.api+json" -d'
{
  "data": {
    "type": "accounts",
    "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
    "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "bic": "NWBKGB22"
    }
  }
}'
