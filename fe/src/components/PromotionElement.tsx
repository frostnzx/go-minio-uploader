import React, { useEffect, useState } from "react";
import { Divider, Modal } from "antd";
import { DownloadOutlined, FileImageTwoTone } from "@ant-design/icons";
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
    setPromotionsList,
}: {
    promoName: string;
    description: string;
    date: dayjs.Dayjs;
    getPromotionsList: Function;
    promotionsList: Promotion[];
    setPromotionsList: Function;
}) {
    const monthAndYearStr = dayjs(date).format("MMMM YYYY");
    const uploadAtStr = "Upload at " + dayjs(date).format("D/M/YYYY h:mm A");
    const handleDelete = async () => {
        // mutate list in state first (for fast ui), then re-fetch from backend to ensure correct list
        setPromotionsList((prev: Promotion[]) =>
            prev.filter((promotion: Promotion) => promotion.name !== promoName)
        ); // make a change in state first
        const response = await axios.delete(
            `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/image-collection/${promoName}`
        );
        getPromotionsList(); // re-fetch

        console.log(response);
    };
    const handleDownload = async () => {};

    const [isModalOpen, setIsModalOpen] = useState(false);
    const showModal = () => {
        setIsModalOpen(true);
    };
    const [images, setImages] = useState<String[]>([]);
    useEffect(() => {
        // since we can't edit or add any image to the collection we just fetch once
        const fetchImages = async () => {
            const response = await axios.get(
                `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/image-collection/${promoName}/uploads`
            );
            setImages(response.data);
        };
        fetchImages();
    }, []);

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
                <div className="w-1/5 flex flex-col justify-center items-center gap-5 text-red-500 font-extrabold">
                    <Button
                        color="default"
                        variant="solid"
                        icon={<DownloadOutlined />}
                        onClick={handleDownload}
                    >
                        Download CSV
                    </Button>
                    <Button onClick={showModal}>View Uploaded List</Button>
                    <Button danger onClick={handleDelete}>
                        Remove Collection
                    </Button>
                </div>
            </div>
            <Divider style={{ borderColor: "#7cb305", marginTop: "50px" }} />

            <Modal
                title={`Uploaded images for ${promoName} collection`}
                closable={{ "aria-label": "Custom Close Button" }}
                open={isModalOpen}
                onCancel={() => setIsModalOpen(false)}
                footer={null}
            >
                <b>
                Total : {images.length} 
                </b>
                <br />
                {images.map((image) => (
                    <div>{image}</div>
                ))}
            </Modal>
        </div>
    );
}
