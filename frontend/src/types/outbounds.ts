import { oTls } from "./outTls"
import { oMultiplex } from "./multiplex"
import { Transport } from "./transport"

export const OutTypes = {
  Direct: 'direct',
  Block: 'block',
  SOCKS: 'socks',
  HTTP: 'http',
  Shadowsocks: 'shadowsocks',
  VMess: 'vmess',
  Trojan: 'trojan',
  Wireguard: 'wireguard',
  Hysteria: 'hysteria',
  VLESS: 'vless',
  ShadowTLS: 'shadowtls',
  TUIC: 'tuic',
  Hysteria2: 'hysteria2',
  Tor: 'tor',
  SSH: 'ssh',
  DNS: 'dns',
  Selector: 'selector',
  URLTest: 'urltest',
}

type OutType = typeof OutTypes[keyof typeof OutTypes]

export interface Dial {
  detour?: string
  bind_interface?: string
  inet4_bind_address?: string
  inet6_bind_address?: string
  routing_mark?: number
  reuse_addr?: boolean
  connect_timeout?: string
  tcp_fast_open?: boolean
  tcp_multi_path?: boolean
  udp_fragment?: boolean
  domain_strategy?: string
  fallback_delay?: string
}

interface OutboundBasics {
  type: OutType
  tag: string
}

export interface WgPeer {
  server: string
  server_port: number
  public_key: string
  pre_shared_key?: string
  allowed_ips?: string[]
  reserved?: number[]
}

export interface Direct extends OutboundBasics, Dial {
  override_address?: string
  override_port?: number
  proxy_protocol?: 0 | 1 | 2
}

export interface Block extends OutboundBasics {}

export interface SOCKS extends OutboundBasics, Dial {
  server: string
  server_port: number
  version?: "4" | "4a" | "5"
  username?: string
  password?: string
  network?: "udp" | "tcp"
  udp_over_tcp?: false | {
    enabled: true
    version?: number
  }
}

export interface HTTP extends OutboundBasics, Dial {
  server: string
  server_port: number
  username?: string
  password?: string
  path?: string
  headers?: {
    [key: string]: string
  }
  tls?: oTls
}

export interface Shadowsocks extends OutboundBasics, Dial {
  server: string
  server_port: number
  method: string
  password: string
  network?: "udp" | "tcp"
  udp_over_tcp?: false | {
    enabled: true
    version?: number
  }
  multiplex?: oMultiplex
}

export interface VMESS extends OutboundBasics, Dial {
  server: string
  server_port: number
  uuid: string
  security?: string
  alter_id: 0
  global_padding?: boolean
  authenticated_length?: boolean
  network?: "udp" | "tcp"
  packet_encoding?: string
  tls?: oTls
  multiplex?: oMultiplex
  transport?: Transport
}

export interface Trojan extends OutboundBasics, Dial {
  server: string
  server_port: number
  password: string
  network?: "udp" | "tcp"
  tls?: oTls
  multiplex?: oMultiplex
  transport?: Transport
}

export interface WireGuard extends OutboundBasics, Dial {
  server?: string
  server_port?: number
  system_interface?: boolean
  gso?: boolean
  interface_name?: string
  local_address: string[]
  private_key: string
  peers?: WgPeer[]
  peer_public_key?: string
  pre_shared_key?: string
  reserved?: number[]
  workers?: number
  mtu?: number
  network?: "udp" | "tcp"
}

export interface Hysteria extends OutboundBasics, Dial {
  server: string
  server_port: number
  up_mbps: number
  down_mbps: number
  obfs?: string
  auth_str?: string
  recv_window_conn?: number
  recv_window?: number
  disable_mtu_discovery?: boolean
  network?: "udp" | "tcp"
  tls: oTls
}

export interface ShadowTLS extends OutboundBasics, Dial {
  server: string
  server_port: number
  version: 1|2|3
  password?: string
  tls: oTls
}

export interface VLESS extends OutboundBasics, Dial {
  server: string
  server_port: number
  uuid: string
  flow?: string
  network?: "udp" | "tcp"
  packet_encoding?: string
  tls?: oTls
  multiplex?: oMultiplex
  transport?: Transport
}

export interface TUIC extends OutboundBasics, Dial {
  server: string
  server_port: number
  uuid: string
  password?: string
  congestion_control?: "cubic"|"new_reno"|"bbr"
  udp_relay_mode?: "native" | "quic"
  udp_over_stream?: boolean
  zero_rtt_handshake?: boolean
  heartbeat?: string
  network?: "udp" | "tcp"
  tls: oTls
}

export interface Hysteria2 extends OutboundBasics, Dial {
  server: string
  server_port: number
  up_mbps?: number
  down_mbps?: number
  obfs?: {
    type?: "salamander"
    password: string
  }
  password?: string
  network?: "udp" | "tcp"
  tls: oTls
  brutal_debug?: boolean
}

export interface Tor extends OutboundBasics, Dial {
  executable_path?: string
  extra_args?: string[]
  data_directory: string
  torrc?: {
    ClientOnly: 0 | 1
  }
}

export interface SSH extends OutboundBasics, Dial  {
  server: string
  server_port?: number
  user?: string
  password?: string
  private_key?: string
  private_key_path?: string
  private_key_passphrase?: string
  host_key?: string[]
  host_key_algorithms?: string[]
  client_version?: string
}

export interface DNS extends OutboundBasics {}

export interface Selector extends OutboundBasics {
  outbounds: string[]
  url?: string
  interval?: string
  tolerance?: number
  idle_timeout?: string
  interrupt_exist_connections?: boolean
}

export interface URLTest extends OutboundBasics {
  outbounds: string[]
  default?: string
  interrupt_exist_connections?: boolean
}

// Create interfaces dynamically based on OutTypes keys
type InterfaceMap = {
  [Key in keyof typeof OutTypes]: {
    type: string
    [otherProperties: string]: any; // You can add other properties as needed
  }
}

// Create union type from InterfaceMap
export type Outbound = InterfaceMap[keyof InterfaceMap]

// Create defaultValues object dynamically
const defaultValues: Record<OutType, Outbound> = {
  direct: { type: OutTypes.Direct },
  block: { type: OutTypes.Block },
  socks: { type: OutTypes.SOCKS, version: "5" },
  http: { type: OutTypes.HTTP, tls: {} },
  shadowsocks: { type: OutTypes.Shadowsocks, method: 'none', multiplex: {} },
  vmess: { type: OutTypes.VMess, tls: {}, multiplex: {}, transport: {}, security: 'auto', global_padding: false },
  trojan: { type: OutTypes.Trojan, tls: {}, multiplex: {}, transport: {} },
  wireguard: { type: OutTypes.Wireguard, local_address: ['10.0.0.2/32','fe80::2/128'], private_key: '' },
  hysteria: { type: OutTypes.Hysteria, up_mbps: 100, down_mbps: 100, tls: { enabled: true } },
  shadowtls: { type: OutTypes.ShadowTLS, version: 3, tls: { enabled: true } },
  vless: { type: OutTypes.VLESS, tls: {}, multiplex: {}, transport: {} },
  tuic: { type: OutTypes.TUIC, congestion_control: 'cubic', tls: { enabled: true } },
  hysteria2: { type: OutTypes.Hysteria2, tls: { enabled: true } },
  tor: { type: OutTypes.Tor, executable_path: './tor', data_directory: '$HOME/.cache/tor', torrc: { ClientOnly: 1 } },
  ssh: { type: OutTypes.SSH },
  dns: { type: OutTypes.DNS },
  selector: { type: OutTypes.Selector },
  urltest: { type: OutTypes.URLTest },
}

export function createOutbound<T extends Outbound>(type: string,json?: Partial<T>): Outbound {
  const defaultObject: Outbound = { ...defaultValues[type], ...(json || {}) }
  return defaultObject
}