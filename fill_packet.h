#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <arpa/inet.h>

typedef struct __attribute__((packed))
{
    char dest[6];
    char sender[6];
    uint16_t protocolType;
} ethhdr;

typedef struct __attribute__((packed))
{
    uint16_t hwType;
    uint16_t protoType;
    char hwLen;
    char protocolLen;
    uint16_t oper;
    char SHA[6];
    char SPA[4];
    char THA[6];
    char TPA[4];
} arphdr;

typedef struct __attribute__((packed))
{
    ethhdr eth;
    arphdr arp;
} arp_packet;

char* fill_arp_packet(char* smac, char* sip);
