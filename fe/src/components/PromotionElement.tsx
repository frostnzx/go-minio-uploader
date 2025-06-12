import React from "react";
import { Divider } from "antd";
import { FileImageTwoTone } from "@ant-design/icons";
import { Button, Flex } from "antd";
import dayjs from "dayjs";
import axios from "axios";
import { Promotion } from "@/interfaces";

export default function PromotionElement({
    promoName,
    description,
    date,
    getPromotionsList,
    promotionsList,
    setPromotionsList
}: {
    promoName: string;
    description: string;
    date: dayjs.Dayjs;
    getPromotionsList : Function;
    promotionsList : Promotion[];
    setPromotionsList : Function;
}) {
    const monthAndYearStr = dayjs(date).format("MMMM YYYY");
    const uploadAtStr = "Upload at " + dayjs(date).format("D/M/YYYY h:mm A");
    const handleDelete = async () => {
        // mutate list in state first (for fast ui), then re-fetch from backend to ensure correct list
        setPromotionsList((prev : Promotion[]) => (
            prev.filter((promotion : Promotion) => promotion.name !== promoName)
        )) // make a change in state first
        const response = await axios.delete(
            `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/image-collection/${promoName}`
        );
        getPromotionsList() // re-fetch

        console.log(response)
    }
    return (
        <div className="w-3/4">
            <div className="flex">
                <div className="w-1/5 flex justify-center items-center">
                    <FileImageTwoTone style={{ fontSize: "50px" }} />
                </div>
                <div className="w-2/5">
                    <div className="text-2xl font-extrabold">{promoName}</div>
                    <div>{monthAndYearStr}</div>
                    <div className="inline-block mt-4 w-auto">
                        Description : {description}
                    </div>
                </div>
                <div className="w-1/5 flex justify-center items-center">
                    {uploadAtStr}
                </div>
                <div className="w-1/5 flex justify-center items-center text-red-500 font-extrabold">
                    <Button danger onClick={handleDelete}>Remove Collection</Button>
                </div>
            </div>
            <Divider style={{ borderColor: "#7cb305", marginTop: "50px" }} />
        </div>
    );
}
