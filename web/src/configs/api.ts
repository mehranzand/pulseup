const API_BASE_URL = "http://localhost:7070/api/"

export const API_URLS = {
    containers_url: `${API_BASE_URL}:host/containers`,
    logs_stream_url: `${API_BASE_URL}logs/stream/localhost/:id`,
    container_event_stream_url: `${API_BASE_URL}events/stream/:host`,
}