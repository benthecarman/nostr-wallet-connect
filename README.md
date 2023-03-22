# nostr-wallet-connect

This service allows you to control your Lightning node or wallet over Nostr.

## Configuration parameters
For self-hosting with your own node, set the following parameters:

- `NOSTR_PRIVKEY`: the private key of this service. Should be a securely randomly generated 32 byte hex string.
- `CLIENT_NOSTR_PUBKEY`: if set, this service will only listen to events authored by this public key. You can set this to your own nostr public key.
- `RELAY`: eg. `wss://relay.damus.io`
- `LN_BACKEND_TYPE`: should be set to `LND`
- `LND_ADDRESS`: the LND gRPC address, eg. `localhost:10009`
- `LND_CERT_FILE"`: the location where LND's `tls.cert` file can be found
- `LND_MACAROON_FILE"`:  the location where LND's `admin.macaroon` file can be found