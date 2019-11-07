#include "fill_packet.h"
#include <stdio.h>

char* fill_arp_packet(char* smac, char* sip)
{
    arp_packet * packet = malloc(sizeof(arp_packet));
    memset(packet, 0, sizeof(arp_packet));
    // Ethernet header
    // Dest = Broadcast (ff:ff:ff:ff:ff)
    memset(packet->eth.dest,0xff,sizeof(packet->eth.dest));


    sscanf(smac,"%hhx:%hhx:%hhx:%hhx:%hhx:%hhx",&packet->eth.sender[0],&packet->eth.sender[1],&packet->eth.sender[2],&packet->eth.sender[3],&packet->eth.sender[4],&packet->eth.sender[5]);

    packet->eth.protocolType = htons(0x0806); // ARP

    // ARP Packet fields
    packet->arp.hwType = htons(1); // Ethernet
    packet->arp.protoType = htons(0x800); //IP;
    packet->arp.hwLen = 6;
    packet->arp.protocolLen = 4;
    packet->arp.oper = htons(2); // response

    // Sender MAC (same as that in eth header)
    memcpy(packet->arp.SHA, packet->eth.sender, 6);

    // Sender IP
    sscanf(sip,"%hhu.%hhu.%hhu.%hhu",&packet->arp.SPA[0],&packet->arp.SPA[1],&packet->arp.SPA[2],&packet->arp.SPA[3]);

    // Dest MAC: Same as SHA, as we use an ARP response
    //memcpy(packet->arp.THA, packet->arp.SHA, 6);

    // Dest IP: Same as SPA
    memcpy(packet->arp.TPA, packet->arp.SPA, 4);

    return (char*) packet;
}
