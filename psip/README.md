TVCT
==== 

Property                      |         Number of Bits | Format       | Byte  Offset
------------------------------|------------------------|--------------|------------
table_id                      |                      8 | 0xC8         | 0
section_syntax_indicator      |                      1 | '1'          | 1
private_indicator             |                      1 | '1'          | 1
reserved                      |                      2 | '11'         | 1
section_length                |                     12 | uimsbf       | 1-2
transport_stream_id           |                     16 | uimsbf       | 3-4
reserved                      |                      2 | '11'         | 5
version_number                |                      5 | uimsbf       | 5
current_next_indicator        |                      1 | bslbf        | 5
section_number                |                      8 | uimsbf       | 6
last_section_number           |                      8 | uimsbf       | 7
protocol_version              |                      8 | uimsbf       | 8 
num_channels_in_section       |                      8 | uimsbf       | 9
channels                      |    N * sizeof(channel) | []channel    | 10
reserved                      |                      6 | '111111'     | 10 + sizeof(channels)
additional_descriptors_length |                     10 | uimsbf       | 10 + sizeof(channels)
additional_descriptors        | N * sizeof(descriptor) | []descriptor | 12 + sizeof(channels)
CRC_32                        |                     32 | rpchof       | section_length - 4


Channel
=======

Property                      |         Number of Bits | Format       | Byte  Offset
------------------------------|------------------------|--------------|------------
short_name                    |                    112 | uimsbf       | 0-13
reserved                      |                      4 | '1111'       | 14
major_channel_number          |                     10 | uimsbf       | 14-15
minor_channel_number          |                     10 | uimsbf       | 15-16
modulation_mode               |                      8 | uimsbf       | 17
carrier_frequency             |                     32 | uimsbf       | 18-21
channel_TSID                  |                     16 | uimsbf       | 22-23
program_number                |                     16 | uimsbf       | 24-25
ETM_location                  |                      2 | uimsbf       | 26
access_controlled             |                      1 | bslbf        | 26
hidden                        |                      1 | bslbf        | 26
reserved                      |                      2 | '11'         | 26
hide_guide                    |                      1 | bslbf        | 26
reserved                      |                      3 | '111'        | 26-27
service_type                  |                      6 | uimsbf       | 27
source_id                     |                     16 | uimsbf       | 28-29
reserved                      |                      6 | '111111'     | 30
descriptors_length            |                     10 | uimsbf       | 30-31
descriptors                   | N * sizeof(descriptor) | []descriptor | 32

