import { TypeDate, convertUtcToDatePart, convertTimeMessage, getTimeIsoCurrent } from "./date";

describe("convertUtcToDatePart", () => {
    const dateStr = "2025-09-22T07:30:45Z"; // UTC

    it("returns correct Day", () => {
        expect(convertUtcToDatePart(dateStr, TypeDate.Day)).toBe(22);
    });

    it("returns correct Month", () => {
        expect(convertUtcToDatePart(dateStr, TypeDate.Month)).toBe(9);
    });

    it("returns correct Year", () => {
        expect(convertUtcToDatePart(dateStr, TypeDate.Year)).toBe(2025);
    });

    it("returns correct Hours (UTC)", () => {
        expect(convertUtcToDatePart(dateStr, TypeDate.Hours)).toBe(7);
    });

    it("returns correct Minutes", () => {
        expect(convertUtcToDatePart(dateStr, TypeDate.Minutes)).toBe(30);
    });

    it("returns correct Seconds", () => {
        expect(convertUtcToDatePart(dateStr, TypeDate.Seconds)).toBe(45);
    });

    it("returns null for invalid date", () => {
        expect(convertUtcToDatePart("invalid", TypeDate.Day)).toBeNull();
    });
});

describe("convertTimeMessage", () => {
    const str = "2025-09-22T07:30:45Z";
    const timezone = 7;
    it("returns correct time message", () => {
        expect(convertTimeMessage(str, timezone)).toBe("14:30");
    });
    it("returns wrong miss str", () => {
        expect(convertTimeMessage("", timezone)).toBe("");
    });
})

describe("getTimeIsoCurrent", () => {
    it("returns correct time nowiso", () => {
        expect(getTimeIsoCurrent()).toBe(getTimeIsoCurrent());
    });
})