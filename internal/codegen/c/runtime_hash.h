// Private streaming SipHash-1-3 implementation. No Telia wrapping arithmetic.
// Native-only deterministic collision seam; not a public hasher callback.
#ifndef _TELIA_KEY_HASH
#define _TELIA_KEY_HASH(value) (value)
#endif
_Static_assert(sizeof(uint64_t) == 8, "SipHash requires 64-bit uint64_t");
typedef struct { uint64_t v0, v1, v2, v3, tail, length; unsigned used; } _TeliaSip;
static inline uint64_t _telia_rotl64(uint64_t x, unsigned n) {
    return (x << n) | (x >> (64 - n));
}
static inline void _telia_sip_round(_TeliaSip *s) {
    s->v0 += s->v1; s->v1 = _telia_rotl64(s->v1,13); s->v1 ^= s->v0; s->v0 = _telia_rotl64(s->v0,32);
    s->v2 += s->v3; s->v3 = _telia_rotl64(s->v3,16); s->v3 ^= s->v2;
    s->v0 += s->v3; s->v3 = _telia_rotl64(s->v3,21); s->v3 ^= s->v0;
    s->v2 += s->v1; s->v1 = _telia_rotl64(s->v1,17); s->v1 ^= s->v2; s->v2 = _telia_rotl64(s->v2,32);
}
static inline _TeliaSip _telia_sip_init(uint64_t k0, uint64_t k1) {
    return (_TeliaSip){.v0=k0^UINT64_C(0x736f6d6570736575), .v1=k1^UINT64_C(0x646f72616e646f6d),
        .v2=k0^UINT64_C(0x6c7967656e657261), .v3=k1^UINT64_C(0x7465646279746573)};
}
static inline void _telia_sip_block(_TeliaSip *s, uint64_t m) {
    s->v3 ^= m; _telia_sip_round(s); s->v0 ^= m;
}
static inline void _telia_sip_write(_TeliaSip *s, const unsigned char *bytes, size_t n) {
    s->length += (uint64_t)n;
    for(size_t i=0; i<n; ++i) {
        s->tail |= (uint64_t)bytes[i] << (8*s->used);
        if (++s->used == 8) { _telia_sip_block(s,s->tail); s->used=0; s->tail=0; }
    }
}
static inline uint64_t _telia_sip_finish(_TeliaSip s) {
    _telia_sip_block(&s,s.tail | ((s.length & 255) << 56));
    s.v2 ^= 255;
    _telia_sip_round(&s); _telia_sip_round(&s); _telia_sip_round(&s);
    return s.v0 ^ s.v1 ^ s.v2 ^ s.v3;
}
