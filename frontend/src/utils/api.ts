import { authApi } from "../api/authApi";
import { toast } from "react-toastify";
import useStore from "../store";

export const sendToTelegram = async ({
  user_id,
  telegram_chat_id,
}: {
  user_id: string;
  telegram_chat_id: "808757394";
}) => {
  const store = useStore();

  try {
    store.setRequestLoading(true);
    const response = await authApi.post("/auth/send-otp-telegram", {
      user_id,
      telegram_chat_id,
    });

    store.setRequestLoading(false);

    if (response.status === 200) {
      console.log("Send to Telegram successfully");
      toast.success("OTP Sent to Telegram Successfully", {
        position: "top-right",
      });
    } else {
      console.error("Failed to send OTP to Telegram");
      toast.error("Failed to send OTP to Telegram", {
        position: "top-right",
      });
    }
  } catch (error: any) {
    console.error("An error occurred:", error);
    store.setRequestLoading(false);
    const errorMessage =
      (error.response &&
        error.response.data &&
        error.response.data.message) ||
      error.response.data.detail ||
      error.message ||
      error.toString();
    toast.error(errorMessage, {
      position: "top-right",
    });
  }
};
