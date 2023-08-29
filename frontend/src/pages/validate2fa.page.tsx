import { object, string, TypeOf } from "zod";
import { useEffect } from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoadingButton } from "../components/LoadingButton";
import { toast } from "react-toastify";
import { Link, useNavigate } from "react-router-dom";
import useStore from "../store";
import { authApi } from "../api/authApi";

const styles = {
  inputField: `form-control block w-full px-4 py-4 text-sm text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none`,
};

const validate2faSchema = object({
  token: string().min(1, "Authentication code is required"),
});

export type Validate2faInput = TypeOf<typeof validate2faSchema>;

const Validate2faPage = () => {
  const navigate = useNavigate();
  const store = useStore();
  const user = store.authUser;

  const {       
    handleSubmit,
    setFocus,
    register,
    formState: { errors },
  } = useForm<Validate2faInput>({
    resolver: zodResolver(validate2faSchema),
  });

  const validate2fa = async (token: string) => {
    try {
      store.setRequestLoading(true);
      const {
        data: { otp_valid },
      } = await authApi.post<{ otp_valid: boolean }>("/auth/otp/validate", {
        token,
        user_id: store.authUser?.id,
      });
      store.setRequestLoading(false);
      if (otp_valid) {
        navigate("/profile");
      } else {
        navigate("/login");
      }
    } catch (error: any) {
      store.setRequestLoading(false);
      const resMessage =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.response.data.detail ||
        error.message ||
        error.toString();
      toast.error(resMessage, {
        position: "top-right",
      });
    }
  };

  const sendToTelegram = async ({
    user_id,
    telegram_chat_id,
  }: {
    user_id: string;
    telegram_chat_id: string;
  }) => {
    try {
      store.setRequestLoading(true);
      const response = await authApi.post("/auth/send-otp-telegram", {
        user_id,
        telegram_chat_id
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

  const handleSendOTPToTelegram = async (
    event: React.MouseEvent<HTMLAnchorElement>
  ) => {
    event.preventDefault(); // Prevent the link from navigating
    console.log("Sending OTP to Telegram");
    
    const user = store.authUser; // Retrieve the user data from the store or elsewhere

    if (!user) {
      console.error("User data is missing");
      return;
    }

    try {
      await sendToTelegram({
        user_id: user.id,
        telegram_chat_id: "808757394",
      });
    } catch (error) {
      console.error("Failed to send OTP to Telegram:", error);
      toast.error("Failed to send OTP to Telegram", {
        position: "top-right",
      });
    }
  };

  const onSubmitHandler: SubmitHandler<Validate2faInput> = (values) => {
    validate2fa(values.token);
  };

  useEffect(() => {
    setFocus("token");
  }, [setFocus]);

  useEffect(() => {
    if (!store.authUser) {
      navigate("/login");
    }
  }, []);

  return (
    <section className="bg-ct-blue-600 min-h-screen grid place-items-center">
      <div className="w-full">
        <h1 className="text-4xl lg:text-6xl text-center font-[600] text-ct-yellow-600 mb-4">
          Welcome Back
        </h1>
        <h2 className="text-lg text-center mb-4 text-ct-dark-200">
          Verify the Authentication Code
        </h2>
        <form
          onSubmit={handleSubmit(onSubmitHandler)}
          className="max-w-md w-full mx-auto overflow-hidden shadow-lg bg-ct-dark-200 rounded-2xl p-8 space-y-5"
        >
          <h2 className="text-center text-3xl font-semibold text-[#142149]">
            Two-Factor Authentication
          </h2>
          <p className="text-center text-sm">
            Open the two-step verification app on your mobile device to get your
            verification code.
          </p>
          <input
            {...register("token")}
            className={styles.inputField}
            placeholder="Authentication Code"
          />
          <p className="mt-2 text-xs text-red-600">
            {errors.token ? errors.token.message : null}
          </p>
  
          <LoadingButton
            loading={store.requestLoading}
            textColor="text-ct-blue-600"
          >
            Authenticate
          </LoadingButton>
  
          <div className="flex flex-col space-y-2 text-center">
            <Link to="/login" className="text-ct-blue-600">
              Back to basic login
            </Link>
            <Link
              to="#"
              className="text-ct-blue-600"
              onClick={handleSendOTPToTelegram}
            >
              Get a new OTP code
            </Link>
          </div>
        </form>
      </div>
    </section>
  );
}

export default Validate2faPage;

