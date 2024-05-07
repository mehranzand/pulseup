
const API_BASE_URL = `${window.location.origin}/api/`

export const API_URLS = {
    containers_url: `${API_BASE_URL}:host/containers`,
    logs_stream_url: `${API_BASE_URL}logs/stream/:host/:id`,
    container_event_stream_url: `${API_BASE_URL}events/stream/:host`,
}