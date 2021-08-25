# Tips

String literals were chosen from the reserved characters of urls as noted [here](https://developers.google.com/maps/documentation/urls/url-encoding#special-characters).

## Search Version 2

From a query with values `foo, bar, baz` encoded as base64:

```text
(((key=company_name:value=BASE64(foo):value_type=STRING:options=CASE_INSENSITIVE_STRING,PARTIAL_MATCH_STRING)*(key=branch_name:value=BASE64(bar):value_type=STRING:options=CASE_INSENSITIVE_STRING,PARTIAL_MATCH_STRING)*(key=trading_name:value=BASE64(baz):value_type=STRING:options=CASE_INSENSITIVE_STRING,PARTIAL_MATCH_STRING))+(key=head_org_id:options=IS_NULL))
```

Resulting filters:

```json
{
  "LinkedFilters": [
    {
      "Filter": null,
      "Type": "OPEN_BRACKET"
    },
    {
      "Filter": null,
      "Type": "OPEN_BRACKET"
    },
    {
      "Filter": null,
      "Type": "OPEN_BRACKET"
    },
    {
      "Filter": {
        "case_insensitive_string": true,
        "key": "company_name",
        "partial_match_string": true,
        "value": "foo"
      },
      "Type": null
    },
    {
      "Filter": null,
      "Type": "CLOSE_BRACKET"
    },
    {
      "Filter": null,
      "Type": "OR"
    },
    {
      "Filter": null,
      "Type": "OPEN_BRACKET"
    },
    {
      "Filter": {
        "case_insensitive_string": true,
        "key": "branch_name",
        "partial_match_string": true,
        "value": "bar"
      },
      "Type": null
    },
    {
      "Filter": null,
      "Type": "CLOSE_BRACKET"
    },
    {
      "Filter": null,
      "Type": "OR"
    },
    {
      "Filter": null,
      "Type": "OPEN_BRACKET"
    },
    {
      "Filter": {
        "case_insensitive_string": true,
        "key": "trading_name",
        "partial_match_string": true,
        "value": "baz"
      },
      "Type": null
    },
    {
      "Filter": null,
      "Type": "CLOSE_BRACKET"
    },
    {
      "Filter": null,
      "Type": "CLOSE_BRACKET"
    },
    {
      "Filter": null,
      "Type": "AND"
    },
    {
      "Filter": null,
      "Type": "OPEN_BRACKET"
    },
    {
      "Filter": {
        "is_null": true,
        "key": "head_org_id"
      },
      "Type": null
    },
    {
      "Filter": null,
      "Type": "CLOSE_BRACKET"
    },
    {
      "Filter": null,
      "Type": "CLOSE_BRACKET"
    },
    {
      "Filter": null,
      "Type": "AND"
    },
    {
      "Filter": {
        "key": "test",
        "value": true
      },
      "Type": null
    }
  ]
}
```
