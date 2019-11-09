package configs

const AppEndpoint = "http://localhost:3000/"
const ChatsRoute = "/applications/{access_token}/chats"
const MessagesRoute = "/applications/{access_token}/chats/{chat_number}/messages"

const RedisChatQueue = "chat"
const ChatWorkerClass = "ChatWorker"
const RedisChatKeyPrefix = "CHAT_"

const MsgRedisQueue = "message"
const MsgWorkerClass = "MessageWorker"
const MsgRedisKeyPrefix = "MSG_"

const RedisAddress = "localhost:6379"
