"use client";
import PromotionElement from "@/components/PromotionElement";
import Image from "next/image";
import {
    Checkbox,
    Divider,
    Form,
    FormProps,
    Input,
    Modal,
    Typography,
} from "antd";
const { Title } = Typography;
import { Button } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { FloatButton } from "antd";
import { useEffect, useReducer, useState } from "react";
import InputForm from "@/components/InputForm";
import dayjs from "dayjs";
import axios from "axios";
import { Promotion } from "@/interfaces";

let formInit: Promotion = {
    name: "",
    description: "",
    date: dayjs(),
};

export default function Home() {
    // Form data
    const [formData, setFormData] = useState(formInit);
    const [fileList, setFileList] = useState<File[]>([]); // also in the form but seperate
    // modal logic
    const [isModalOpen, setIsModalOpen] = useState(false);
    const showModal = () => {
        setIsModalOpen(true);
    };

    // fetching logic
    const [promotionsList, setPromotionsList] = useState<Promotion[]>([]);
    const getPromotionsList = async () => {
        const response = await axios.get(
            `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/image-collections`
        );
        setPromotionsList(response.data);
    };
    // initial fetch
    useEffect(() => {
        getPromotionsList();
    }, []);

    // Submit handler send into form prop
    const submitHandler = async () => {
        console.log(formData);
        console.log(fileList);
        setIsModalOpen(false);

        // then fetch to backend to create new promotion
        const data = new FormData();
        const jsonInfoData = {
            name: formData.name,
            description: formData.description,
            date: formData.date.toISOString(),
        };
        data.append("info", JSON.stringify(jsonInfoData));
        fileList.forEach((file) => data.append("images", file));
        try {
            const response = await axios.post(
                `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/image-collection`,
                data
            );
            getPromotionsList(); // re-fetch

            console.log("Success submitting new promotion : ", response.data);
        } catch (err) {
            console.log("Error submitting new promotion : ", err);
        }
    };

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
                    onClick={showModal}
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
                {promotionsList.map((promotion: Promotion) => (
                    <PromotionElement
                        key={promotion.name}
                        promoName={promotion.name}
                        description={promotion.description}
                        date={promotion.date}
                        getPromotionsList={getPromotionsList}
                        promotionsList={promotionsList}
                        setPromotionsList={setPromotionsList}
                    />
                ))}
            </div>

            <Modal
                title="Add Promotion"
                closable={{ "aria-label": "Custom Close Button" }}
                open={isModalOpen}
                onCancel={() => setIsModalOpen(false)}
                footer={null}
            >
                <InputForm
                    onSubmit={submitHandler}
                    setFormData={setFormData}
                    formData={formData}
                    fileList={fileList}
                    setFileList={setFileList}
                />
            </Modal>
        </div>
    );
}
