"""Bounded read-only Delta India public contract evidence collection.

The output is an observation, not a verified lot size, fee or margin rule.
Do not use it to authorize trading.
"""
import datetime
import http.client
import json
import sys
import urllib.parse

HOST = "api.india.delta.exchange"
MAX_BYTES = 4 * 1024 * 1024
CONTRACT_TYPES = ("perpetual_futures", "futures", "call_options", "put_options")
FIELDS = (
    "id", "symbol", "contract_type", "state", "trading_status",
    "notional_type", "is_quanto", "contract_value", "contract_unit_currency",
    "tick_size", "initial_margin", "maintenance_margin", "maker_commission_rate",
    "taker_commission_rate", "settlement_time", "strike_price",
    "position_size_limit", "price_band",
)


def unique_object(items):
    result = {}
    for key, value in items:
        if key in result:
            raise ValueError("duplicate response key")
        result[key] = value
    return result


def inspect_type(contract_type):
    query = urllib.parse.urlencode({
        "contract_types": contract_type,
        "states": "live",
        "page_size": "2",
    })
    conn = http.client.HTTPSConnection(HOST, timeout=12)
    try:
        conn.request("GET", "/v2/products?" + query, headers={
            "Accept": "application/json", "Accept-Encoding": "identity",
            "User-Agent": "Delta-Crypto-Algo-public-contract-probe/1",
        })
        response = conn.getresponse()
        if response.status != 200:
            raise ValueError(f"HTTP {response.status}")
        if response.getheader("Content-Encoding", "identity").lower() != "identity":
            raise ValueError("unexpected response encoding")
        body = response.read(MAX_BYTES + 1)
        if len(body) > MAX_BYTES:
            raise ValueError("response exceeds byte limit")
    finally:
        conn.close()
    payload = json.loads(body, object_pairs_hook=unique_object)
    if not isinstance(payload, dict) or payload.get("success") is not True:
        raise ValueError("unsuccessful response envelope")
    records = payload.get("result")
    if not isinstance(records, list) or len(records) > 2:
        raise ValueError("unexpected result page size")
    meta = payload.get("meta")
    cursor = "missing"
    if isinstance(meta, dict) and "after" in meta:
        cursor = type(meta["after"]).__name__
    print("public_page " + json.dumps({
        "type": contract_type,
        "records": len(records),
        "after_type": cursor,
    }, ensure_ascii=True))
    for record in records:
        if not isinstance(record, dict):
            raise ValueError("non-object product")
        if record.get("contract_type") != contract_type:
            raise ValueError("unexpected product type")
        sample = {key: record[key] for key in FIELDS if key in record}
        for key in ("underlying_asset", "quoting_asset", "settling_asset"):
            obj = record.get(key)
            if isinstance(obj, dict):
                sample[key] = {"symbol": obj.get("symbol")}
        specs = record.get("product_specs")
        if isinstance(specs, dict):
            sample["product_specs"] = {
                key: specs[key] for key in ("only_reduce_only_orders_allowed", "tags")
                if key in specs
            }
        print("public_sample " + json.dumps(sample, ensure_ascii=True))


def main():
    print("public_observed_at " + datetime.datetime.now(datetime.timezone.utc).isoformat())
    failures = 0
    for kind in CONTRACT_TYPES:
        try:
            inspect_type(kind)
        except (OSError, ValueError, http.client.HTTPException) as exc:
            print("public_probe_error " + json.dumps({
                "type": kind, "error": str(exc)[:300],
            }, ensure_ascii=True))
            failures += 1
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
