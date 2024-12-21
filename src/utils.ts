import React from "react";

function humanFileSize(bytes: number, si: boolean = false, dp: number = 2): string {
    const thresh = si ? 1000 : 1024;

    if (Math.abs(bytes) < thresh) {
        return bytes + ' B';
    }

    const units = si
        ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
    let u = -1;
    const r = 10 ** dp;

    do {
        bytes /= thresh;
        ++u;
    } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


    return bytes.toFixed(dp) + ' ' + units[u];
}

function timeSince(date: Date): string {
    var seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    const suffix = seconds < 0 ? ' from now' : ' ago';

    var interval = seconds / 31536000;

    if (interval > 1) {
        return Math.floor(interval) + " years" + suffix;
    }
    interval = seconds / 2592000;
    if (interval > 1) {
        return Math.floor(interval) + " months" + suffix;
    }
    interval = seconds / 86400;
    if (interval > 1) {
        return Math.floor(interval) + " days" + suffix;
    }
    interval = seconds / 3600;
    if (interval > 1) {
        return Math.floor(interval) + " hours" + suffix;
    }
    interval = seconds / 60;
    if (interval > 1) {
        return Math.floor(interval) + " minutes" + suffix;
    }

    return Math.floor(seconds) + " second" + (seconds > 1 ? 's' : '') + suffix;
}

function useInterval(callback: () => void, delay: number | null) {
    const intervalRef = React.useRef<number>();
    const callbackRef = React.useRef<() => void>(callback);

    // Remember the latest callback:
    //
    // Without this, if you change the callback, when setInterval ticks again, it
    // will still call your old callback.
    //
    // If you add `callback` to useEffect's deps, it will work fine but the
    // interval will be reset.

    React.useEffect(() => {
        callbackRef.current = callback;
    }, [callback]);

    // Set up the interval:

    React.useEffect(() => {
        if (typeof delay === 'number') {
            intervalRef.current = window.setInterval(() => callbackRef.current(), delay);

            // Clear interval if the components is unmounted or the delay changes:
            return () => window.clearInterval(intervalRef.current);
        }
    }, [delay]);

    // Returns a ref to the interval ID in case you want to clear it manually:
    return intervalRef;
}

export default { humanFileSize, timeSince, useInterval };
