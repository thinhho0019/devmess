 

export const TypeDate = {
    Day: "day",
    Month: "month",
    Year: "year",
    Hours: "hours",
    Minutes: "minutes",
    Seconds: "seconds",
};

export type TypeDate = (typeof TypeDate)[keyof typeof TypeDate];

export const convertUtcToDatePart = (str: string, type: TypeDate, timezoneOffset = 0): number | null => {
    const d = new Date(str);
    if (isNaN(d.getTime())) return null;
    const localTime = new Date(d.getTime() + timezoneOffset * 60 * 60 * 1000);
    switch (type) {
        case TypeDate.Day:
            return localTime.getUTCDate();
        case TypeDate.Month:
            return localTime.getUTCMonth() + 1;
        case TypeDate.Year:
            return localTime.getUTCFullYear();
        case TypeDate.Hours:
            return localTime.getUTCHours();
        case TypeDate.Minutes:
            return localTime.getUTCMinutes();
        case TypeDate.Seconds:
            return localTime.getUTCSeconds();
        default:
            return null;
    }
};
export const convertTimeToOnlineStatus = (seconds: string ): string => {
    const secondss = Math.floor(Number(seconds));
    console.log("convertTimeToOnlineStatus seconds:", seconds);
    if (isNaN(secondss) || secondss < 0) return "unknown";
    if (secondss < 60) {
        return `${secondss} seconds ago`;
    }
    const minutes = Math.floor(secondss / 60);
    if (minutes < 60) {
        return `${minutes} minutes ago`;
    }
    const hours = Math.floor(minutes / 60);
    if (hours < 24) {
        return `${hours} hours ago`;
    }
    const days = Math.floor(hours / 24);
    return `${days} days ago`;
};
/**
 * Convert UTC time string to formatted local time based on timezone offset.
 * @param str - UTC time string (ví dụ: "2025-09-22T07:30:00Z")
 * @param timezone - múi giờ cần chuyển (ví dụ: +7 cho Việt Nam)
 * @returns string thời gian đã format (HH:mm) hoặc "" nếu input không hợp lệ
 */
export const convertTimeMessage = (str: string, timezone: number): string | "" => {
    if (!str) return "";
    try {
        // Parse chuỗi UTC
        const utcDate = new Date(str);
        if (isNaN(utcDate.getTime())) return "";

        // Lấy timestamp UTC (ms)
        const utcTime = utcDate.getTime();

        // Tính offset
        const localTime = new Date(utcTime + timezone * 60 * 60 * 1000);

        // Format: HH:mm with proper padding
        const pad = (n: number) => String(n).padStart(2, "0");
        const hours = pad(localTime.getUTCHours());
        const minutes = pad(localTime.getUTCMinutes());
        return `${hours}:${minutes}`;
    } catch {
        return "";
    }
};
/**
 * Get time utc current .
 * @returns time utc type string hoặc null if input dont match
 */
export const getTimeIsoCurrent = (): string => {
    const now = new Date();
    const isoNoMs = now.toISOString().split('.')[0] + "Z";
    return isoNoMs
}