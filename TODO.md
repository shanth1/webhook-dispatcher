# TODO

## Improvements

- [x] Refactoring to a modular/plugin architecture
- [ ] Add URL validation for `base_url` in Kanboard adapter configuration.
- [ ] Move `disable_unknown_templates` parameter to individual `WebhookConfig` scope instead of global scope.
- [ ] Add metrics/Prometheus integration (middleware).
- [ ] Add retry mechanism for failed notifications (Outbound adapters).

## Refactoring

- [ ] Refactor Kanboard payload parsing logic (currently relies heavily on interface{} maps).
