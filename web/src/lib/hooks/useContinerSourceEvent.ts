import { useEffect } from "react";
import { API_URLS } from "../../configs/api";
import { useAppDispatch } from "../../hooks"
import { handleSourceEvent } from "../../stores/slices/containerSlice";

const useContinerSourceEvent = () => {
    let es: EventSource | null = null;

    const dispatch = useAppDispatch();

    useEffect(() => {
        es = new EventSource(API_URLS.container_event_stream_url.replace(':host', 'localhost'));
        es.addEventListener("container", (e: Event) => {
            handleEvent((e as MessageEvent));
        });

        function handleEvent(event: MessageEvent) {
            dispatch(handleSourceEvent(JSON.parse(event.data)))
        }
        return () => es?.close();
    }, []);
}

export default useContinerSourceEvent;