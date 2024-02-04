export function formatTimestamp(timestamp) {
    timestamp = parseInt(timestamp);
    let date;
    if (timestamp.toString().length === 13) {
        // 如果长度为 13，则为毫秒级时间戳
        date = new Date(timestamp);
    } else {
        // 否则为普通时间戳
        date = new Date(timestamp * 1000);
    }
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hours = date.getHours();
    const minutes = date.getMinutes();
    const seconds = date.getSeconds();
    function addLeadingZero(value) {
        return value < 10 ? `0${value}` : value;
    }
    return `${year}-${addLeadingZero(month)}-${addLeadingZero(day)},${addLeadingZero(hours)}:${addLeadingZero(minutes)}:${addLeadingZero(seconds)}`;
}


export function formatBytes(bytes) {
    try {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    } catch {
        return "0 Bytes"
    }
}

export function base64Decode(str) {
    try {
        return atob(str);
    } catch (error) {
        console.error("Base64 解码失败:", error);
        return null;
    }
}


// 判断新数据是否与现有数据不同
export function isNewData(newData, existingData) {
    return JSON.stringify(newData) !== JSON.stringify(existingData);
}
