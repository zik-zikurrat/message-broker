Magic - ZAKA
Version
Magic+Version = 5 bytes
Message Types - 5 types(TypePing, TypePong, TypePublish, TypeSubscribe, TypeAck, TypeError)

What need inside message:
    Publish: topic, ack_level, payload
    Subscribe: topic, consumer_group
    Ping: nothing
    Ack: correlation_id
