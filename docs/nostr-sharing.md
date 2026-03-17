# Nostr — partage de listes de feeds

## Contexte

Ratatosk permet de partager sa liste de feeds abonnés via le protocole
[Nostr](https://nostr.com). Nostr est décentralisé : chaque utilisateur a une
keypair secp256k1, publie des événements sur des relais publics (WebSocket), et
n'a pas besoin de compte ni de serveur central.

Le format retenu est **NIP-78, kind 30078** (application-specific replaceable
event) avec `d: "ratatosk-feeds"` et un tag `r` par URL RSS. C'est un
événement *replaceable* : chaque nouvelle publication remplace la précédente
sur les relais, sans accumulation.

---

## Utilisation

```bash
# Générer une keypair (une seule fois)
ratatosk nostr init
# → Public key: npub1...

# Publier ta liste de feeds
ratatosk nostr publish
# → Published feed list to 3/3 relays.

# Importer la liste d'un ami (npub ou hex pubkey)
ratatosk nostr import npub1...
# → Imported 12 feeds (3 already subscribed, 0 errors).
```

---

## Architecture

```
domain          — aucun changement
application     — NostrPort (interface) + NostrService
infrastructure  — nostr/client.go (implémente NostrPort)
cmd             — nostr.go (init | publish | import)
config          — NostrPrivKey, NostrRelays
```

La structure miroir exacte du fetcher RSS existant :
- `GoFeedFetcher` → `NostrClient` (driven adapter dans infrastructure)
- `Fetcher` interface → `NostrPort` interface (dans application)

---

## Fichiers

| Fichier | Rôle |
|---------|------|
| `internal/config/config.go` | +`NostrPrivKey`, `NostrRelays`, `SaveNostrPrivKey()` |
| `internal/application/nostr_port.go` | Interface `NostrPort` |
| `internal/application/nostr_service.go` | `GenerateKey`, `PublishFeedList`, `ImportFeedList` |
| `internal/infrastructure/nostr/client.go` | Implémentation go-nostr |
| `cmd/ratatosk/nostr.go` | Commandes cobra |
| `cmd/ratatosk/main.go` | Wiring |

---

## Interface NostrPort

```go
type NostrPort interface {
    GeneratePrivKey() (string, error)
    DerivePubKey(privKeyHex string) (string, error)
    NormalizePubKey(input string) (string, error) // npub bech32 ou hex brut
    PublishFeedList(ctx context.Context, privKeyHex string, urls []string) error
    FetchFeedList(ctx context.Context, pubKeyHex string) ([]string, error)
}
```

---

## Format de l'événement Nostr

```json
{
  "kind": 30078,
  "tags": [
    ["d", "ratatosk-feeds"],
    ["r", "https://lobste.rs/rss"],
    ["r", "https://example.com/feed.xml"]
  ],
  "content": "",
  "pubkey": "<hex pubkey>",
  "sig": "<schnorr sig>"
}
```

Kind 30078 est *parameterised replaceable* : les relais conservent uniquement
le dernier event par `(pubkey, kind, d-tag)`. Republier met à jour la liste
proprement.

---

## Config

Champs ajoutés dans `~/.config/ratatosk/config.yaml` :

```yaml
nostr_priv_key: ""        # généré par `ratatosk nostr init`
nostr_relays:
  - wss://relay.damus.io
  - wss://nos.lol
  - wss://relay.nostr.band
```

La clé privée est stockée avec les permissions `0600`. Elle n'est jamais
affichée dans les sorties — seulement la clé publique (npub).

---

## Gestion des erreurs

| Cas | Comportement |
|-----|-------------|
| Clé absente au publish | Erreur : "run `ratatosk nostr init` first" |
| Tous les relais injoignables | Erreur fatale |
| Certains relais échouent | Warning + succès partiel |
| npub invalide | Erreur fatale |
| URL invalide dans l'event distant | Skip + warning, continue |
| Feed déjà souscrit | Skip silencieux |

---

## Dépendance

```
github.com/nbd-wtf/go-nostr
```

Pure Go, pas de CGo.

---

## Ordre d'implémentation

1. `go get github.com/nbd-wtf/go-nostr`
2. Étendre `internal/config/config.go`
3. Créer `internal/application/nostr_port.go`
4. Créer `internal/application/nostr_service.go` + tests
5. Créer `internal/infrastructure/nostr/client.go`
6. Créer `cmd/ratatosk/nostr.go` + wiring dans `main.go`
