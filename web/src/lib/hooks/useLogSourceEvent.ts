import { useEffect, useState, useRef } from "react";
import { useLocation, useParams } from 'react-router-dom';
import { API_URLS } from "../../configs/api";
import { Log } from "../../types/Log";
import { useAppSelector } from "../../hooks";
import { Container } from "../../types/Container";
import _ from 'lodash'

export default function useLogSourceEvent() {
    const location = useLocation();
    let params = useParams()
    const container = useAppSelector((state) => state.containers.data.find(a => a.id == params.id)) as Container
    let es: EventSource | null = null
    let debounceBuffer = _.debounce(flushBuffer, 250, { maxWait: 500 })
    let buffer: Log[] = []
    let [messages, _setMessages] = useState<Log[]>([])
    let [loading, setLoading] = useState(true)
    const messagesRef = useRef(messages)
    const pausedRef = useRef(false)
    const maxMessagesArraySize = 400

    const setMessages = (logs: Log[]) => {
        messagesRef.current = logs;
        _setMessages(logs);
    };

    function flushBuffer() {
        if (!pausedRef.current) {
            if (messagesRef.current.length >= maxMessagesArraySize) {
                var sliced = _.slice(messagesRef.current, buffer.length)
                setMessages([...sliced, ...buffer])
            }
            else {
                setMessages([...messagesRef.current, ...buffer])
            }
        }

        setLoading(false)
        buffer = []
    }

    function close() : void {
        es?.close()
    }

    function clear() : void {
        buffer = []
        setMessages([])
        setLoading(true)
        debounceBuffer.cancel()
        pausedRef.current = false
        messagesRef.current = []
    }

    function pause() : boolean {
        if (!pausedRef.current) {
            pausedRef.current = true
        }

        return pausedRef.current
    }

    function resume() : boolean {
        if (pausedRef.current) {
            pausedRef.current = false
        }

        return pausedRef.current
    }

    useEffect(() => {
        if(container === undefined) return

        es = new EventSource(API_URLS.logs_stream_url.replace(':host', container.host).replace(':id', container.id));
        es.onmessage = (e) => {
            handleEvent(e.data);
        };

        es.addEventListener("error", (e) => {
            if (es?.readyState === EventSource.CLOSED) {
                console.log('STREMING CLOSED', e)
            }
        });

        es.addEventListener("container-stopped", () => {
            const stopped: Log = {
                message: 'container-stopped'
            } as Log
            buffer.push(stopped)
            debounceBuffer.flush()
        });

        function handleEvent(data: any) {
            var message = JSON.parse(data)
            const converted: Log = {
                message: message.m,
                date: message.ts as number,
                type: message.t,
                rowType: "log"
            };
            buffer.push(converted)
            debounceBuffer()
        }

        return () => {
            clear()
            close()
        }
    }, [location, container])

    return { messages, loading, pause, resume }
}