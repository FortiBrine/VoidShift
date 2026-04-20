export type NetworkSummary = {
  id: number
  name: string
  address: string
  listen_port: number
}

export type NetworkDetails = {
  id: number
  public_key: string
  address: string
  listen_port: number
  peers: PeerSummary[]
}

export type PeerSummary = {
  id: number
  public_key: string
  allowed_ips: string[]
}

type NetworksResponse = {
  networks: NetworkSummary[]
}

type CreateNetworkRequest = {
  name: string
  address: string
  listen_port: number
}

type CreateNetworkResponse = {
  id: number
  public_key: string
  address: string
  listen_port: number
}

type CreatePeerRequest = {
  allowed_ips: string[]
}

type CreatePeerResponse = {
  id: number
  public_key: string
}

export const useWireguardApi = () => {
  const api = useApiClient()

  const listNetworks = async () => {
    const response = await api.request<NetworksResponse>('/api/vpn/wireguard/networks')
    return response.networks
  }

  const getNetwork = async (id: number) => {
    return await api.request<NetworkDetails>(`/api/vpn/wireguard/networks/${id}`)
  }

  const createNetwork = async (payload: CreateNetworkRequest) => {
    return await api.request<CreateNetworkResponse, CreateNetworkRequest>('/api/vpn/wireguard/networks/generate', {
      method: 'POST',
      body: payload,
    })
  }

  const addPeer = async (networkId: number, payload: CreatePeerRequest) => {
    return await api.request<CreatePeerResponse, CreatePeerRequest>(`/api/vpn/wireguard/networks/${networkId}/peers/generate`, {
      method: 'POST',
      body: payload,
    })
  }

  const getPeerConfig = async (peerId: number) => {
    return await api.request<string>(`/api/vpn/wireguard/peers/${peerId}/config`, {
      responseType: 'text',
    })
  }

  const getPeerQr = async (peerId: number) => {
    return await api.request<Blob>(`/api/vpn/wireguard/peers/${peerId}/qr`, {
      responseType: 'blob',
    })
  }

  const networkUp = async (networkId: number) => {
    await api.request<void>(`/api/vpn/wireguard/networks/${networkId}/up`, {
      method: 'POST',
    })
  }

  const networkDown = async (networkId: number) => {
    await api.request<void>(`/api/vpn/wireguard/networks/${networkId}/down`, {
      method: 'POST',
    })
  }

  return {
    listNetworks,
    getNetwork,
    createNetwork,
    addPeer,
    getPeerConfig,
    getPeerQr,
    networkUp,
    networkDown,
  }
}
