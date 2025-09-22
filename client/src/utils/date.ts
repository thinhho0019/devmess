export const TypeDate = {
    Day: "day",
    Month: "month",
    Year: "year",
    Hours: "hours",
    Minutes: "minutes",
    Seconds: "seconds",
};

export type TypeDate = (typeof TypeDate)[keyof typeof TypeDate];

export const convertUtcToDatePart = (str: string, type: TypeDate): number | null => {
    const d = new Date(str);
    if (isNaN(d.getTime())) return null;

    switch (type) {
        case TypeDate.Day:
            return d.getUTCDate();
        case TypeDate.Month:
            return d.getUTCMonth() + 1;
        case TypeDate.Year:
            return d.getUTCFullYear();
        case TypeDate.Hours:
            return d.getUTCHours();
        case TypeDate.Minutes:
            return d.getUTCMinutes();
        case TypeDate.Seconds:
            return d.getUTCSeconds();
        default:
            return null;
    }
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

        // Format: HH:mm:ss dd/MM/yyyy
        const pad = (n: number) => n.toString().padStart(2, "0");
        const hours = pad(localTime.getUTCHours());
        const minutes = pad(localTime.getUTCMinutes());
        return `${hours}:${minutes}`;
    } catch (e) {
        return "";
    }
};
/**
 * Get time utc current .
 * @returns time utc type string hoặc null if input dont match
 */
export const getTimeIsoCurrent = ():string=>{
    const now = new Date();
    const isoNoMs = now.toISOString().split('.')[0] + "Z";
    return isoNoMs
}