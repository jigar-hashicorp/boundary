# ***
# WARNING: Do not EDIT or MERGE this file, it is generated by packagespec.
# ***

LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_ID             := 00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714
LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_TYPE           := base
LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_BASE_LAYER     := 
LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_SOURCE_INCLUDE := 
LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_SOURCE_EXCLUDE := 
LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_CACHE_KEY_FILE := .buildcache/cache-keys/base-ecc4855907c70055bfa6db2f6c47b9702dd39714
LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_ARCHIVE_FILE   := .buildcache/archives/00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714.tar.gz
$(eval $(call LAYER,$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_ID),$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_TYPE),$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_BASE_LAYER),$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_SOURCE_INCLUDE),$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_SOURCE_EXCLUDE),$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_CACHE_KEY_FILE),$(LAYER_00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714_ARCHIVE_FILE)))

LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_ID             := 01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44
LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_TYPE           := ui
LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_BASE_LAYER     := 00-base-ecc4855907c70055bfa6db2f6c47b9702dd39714
LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_SOURCE_INCLUDE := internal/ui/VERSION
LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_SOURCE_EXCLUDE := 
LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_CACHE_KEY_FILE := .buildcache/cache-keys/ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44
LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_ARCHIVE_FILE   := .buildcache/archives/01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44.tar.gz
$(eval $(call LAYER,$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_ID),$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_TYPE),$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_BASE_LAYER),$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_SOURCE_INCLUDE),$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_SOURCE_EXCLUDE),$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_CACHE_KEY_FILE),$(LAYER_01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44_ARCHIVE_FILE)))

LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_ID             := 02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add
LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_TYPE           := go-modules
LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_BASE_LAYER     := 01-ui-58c5bc2a874511c5f54b529956fe3d3b56ac4f44
LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_SOURCE_INCLUDE := go.mod go.sum */go.mod */go.sum
LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_SOURCE_EXCLUDE := 
LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_CACHE_KEY_FILE := .buildcache/cache-keys/go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add
LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_ARCHIVE_FILE   := .buildcache/archives/02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add.tar.gz
$(eval $(call LAYER,$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_ID),$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_TYPE),$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_BASE_LAYER),$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_SOURCE_INCLUDE),$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_SOURCE_EXCLUDE),$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_CACHE_KEY_FILE),$(LAYER_02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add_ARCHIVE_FILE)))

LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_ID             := 03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8
LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_TYPE           := copy-source
LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_BASE_LAYER     := 02-go-modules-d9a667a5f98d65166aad18ef355d02c6a77c4add
LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_SOURCE_INCLUDE := *.go
LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_SOURCE_EXCLUDE := 
LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_CACHE_KEY_FILE := .buildcache/cache-keys/copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8
LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_ARCHIVE_FILE   := .buildcache/archives/03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8.tar.gz
$(eval $(call LAYER,$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_ID),$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_TYPE),$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_BASE_LAYER),$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_SOURCE_INCLUDE),$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_SOURCE_EXCLUDE),$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_CACHE_KEY_FILE),$(LAYER_03-copy-source-32dc6b3b392d832be3c4fae06f8d8096796eb9f8_ARCHIVE_FILE)))
