import dayjs from "dayjs"

export interface Promotion {
    name: string;
    description: string;
    date: dayjs.Dayjs;
}