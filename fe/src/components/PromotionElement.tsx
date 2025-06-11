import React from "react";
import { Divider } from "antd";
import { FileImageTwoTone } from "@ant-design/icons";
import { Button, Flex } from "antd";

export default function PromotionElement({
    promoName,
    description,
}: {
    promoName: string;
    description: string;
}) {
    // will decide about date format later
    return (
        <div className="w-3/4">
            <div className="flex">
                <div className="w-1/5 flex justify-center items-center">
                    <FileImageTwoTone style={{ fontSize: "50px" }} />
                </div>
                <div className="w-2/5">
                    <div className="text-2xl font-extrabold">{promoName}</div>
                    <div>September 2025</div>
                    <div className="inline-block mt-4 w-auto">
                        Description : {description}
                    </div>
                </div>
                <div className="w-1/5 flex justify-center items-center">
                    Upload at 19/5/2025 12:33 PM
                </div>
                <div className="w-1/5 flex justify-center items-center text-red-500 font-extrabold">
                    <Button danger>Remove Collection</Button>
                </div>
            </div>
            <Divider style={{ borderColor: "#7cb305", marginTop: "50px" }} />
        </div>
    );
}
