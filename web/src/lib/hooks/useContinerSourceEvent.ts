import { useEffect } from "react";
import { API_URLS } from "../../configs/api";
import { useAppDispatch, useAppSelector } from "../../hooks"
import { handleSourceEvent } from "../../stores/slices/containerSlice";

const useContinerSourceEvent = () => {
    let es: EventSource | null = null;
    const { current } = useAppSelector((state) => state.host)
    const dispatch = useAppDispatch();

    useEffect(() => {
        if (current === undefined) return

        es = new EventSource(API_URLS.container_event_stream_url.replace(':host', current));
        es.addEventListener("container", (e: Event) => {
            handleEvent((e as MessageEvent));
        });

        function handleEvent(event: MessageEvent) {
            dispatch(handleSourceEvent(JSON.parse(event.data)))
        }
        return () => es?.close();
    }, [current]);
}

export default useContinerSourceEvent;