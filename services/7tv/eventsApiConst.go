package seventv

import "time"

const EAPI_HEARTBEAT_LEV = 5 * time.Second

const EAPI_OP_HELLO = 1
const EAPI_OP_HEARTBEAT = 2
const EAPI_OP_RECONNECT = 4
const EAPI_OP_SUBSCRIBE = 35

const EAPI_CONDITION_OBJECT_ID = "object_id"
const EAPI_CONDITION_HOST_ID = "host_id"
const EAPI_CONDITION_CONNECTION_ID = "connection_id"
const EAPI_CONDITION_CTX = "ctx"
const EAPI_CONDITION_PLATFORM = "platform"
const EAPI_CONDITION_ID = "id"

const EAPI_CTX_CHANNEL = "channel"
const EAPI_PLATFORM_TWITCH = "TWITCH"

const EAPI_SUB_ENTITLEMENT_CREATE = "entitlement.create"
const EAPI_SUB_ENTITLEMENT_UPDATE = "entitlement.update"
const EAPI_SUB_ENTITLEMENT_DELETE = "entitlement.delete"
const EAPI_SUB_ENTITLEMENT_ALL = "entitlement.*"

const EAPI_SUB_COSMETIC_CREATE = "cosmetic.create"
const EAPI_SUB_COSMETIC_UPDATE = "cosmetic.update"
const EAPI_SUB_COSMETIC_DELETE = "cosmetic.delete"
const EAPI_SUB_COSMETIC_ALL = "cosmetic.*"
