"use client";
import PromotionElement from "@/components/PromotionElement";
import Image from "next/image";
import { Divider, Typography } from "antd";
const { Title } = Typography;
import { Button } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { FloatButton } from "antd";

export default function Home() {
    const handleAdd = () => {};
    return (
        <div className="">
            <Title level={2} style={{ margin: "50px" }}>
                Promotion Comparison
                <div></div>
            </Title>
            <div className="flex justify-center items-center my-[50px]">
                <Button
                    size="large"
                    variant="outlined"
                    icon={<PlusOutlined />}
                    onClick={handleAdd}
                    style={{
                        display: "flex",
                        alignItems: "center",
                        gap: 8,
                        fontWeight: "600",
                        borderRadius: 8,
                        padding: "0 20px",
                    }}
                >
                    Add Collection
                </Button>
            </div>
            <div className="flex flex-col justify-center items-center">
                <PromotionElement
                    promoName={"BigC"}
                    description={"Very nice"}
                />
                <PromotionElement
                    promoName={"BigC"}
                    description={"Very nice"}
                />
                <PromotionElement
                    promoName={"BigC"}
                    description={"Very nice"}
                />

            </div>
        </div>
    );
}
