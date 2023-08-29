import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import { authApi } from "../api/authApi";
import { IUser } from "../api/types";
import TwoFactorAuth from "../components/TwoFactorAuth";
import useStore from "../store";


const ProfilePage = () => {
  const [secret, setSecret] = useState({
    otp_token: ""
  });
  const [openModal, setOpenModal] = useState(false);
  const navigate = useNavigate();
  const store = useStore();
  const user = store.authUser;

  const generateOTPCode = async ({
    user_id,
  }: {
    user_id: string;
  }) => {
    try {
      store.setRequestLoading(true);
      const response = await authApi.post<{
        otp_token: string;
      }>("/auth/otp/generate", { user_id});
      store.setRequestLoading(false);

      if (response.status === 200) {
        setOpenModal(true);
        console.log({
          otp_token: response.data.otp_token,
        });
        setSecret({
          otp_token: response.data.otp_token,
        });
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

  const disableTwoFactorAuth = async (user_id: string) => {
    try {
      store.setRequestLoading(true);
      const {
        data: { user },
      } = await authApi.post<{
        otp_disabled: boolean;
        user: IUser;
      }>("/auth/otp/disable", { user_id });
      store.setRequestLoading(false);
      store.setAuthUser(user);
      toast.warning("Two Factor Authentication Disabled", {
        position: "top-right",
      });
    } catch (error: any) {
      store.setRequestLoading(false);
      const resMessage =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.message ||
        error.toString();
      toast.error(resMessage, {
        position: "top-right",
      });
    }
  };

  useEffect(() => {
    console.log(store.authUser);
    if (!store.authUser) {
      navigate("/login");
    }
  }, []);

  return (
    <>
      <section className="bg-ct-blue-600  min-h-screen pt-10">
        <div className="max-w-4xl p-12 mx-auto bg-ct-dark-100 rounded-md h-[20rem] flex gap-20 justify-center items-start">
          <div className="flex-grow-2">
            <h1 className="text-2xl font-semibold">Profile Page</h1>
            <div className="mt-8">
              <p className="mb-4">ID: {user?.id}</p>
              <p className="mb-4">Name: {user?.name}</p>
              <p className="mb-4">Email: {user?.email}</p>
              <p className="mb-4">Chat Telegram ID: 808757394</p>
            </div>
          </div>
          <div>
            <h3 className="text-2xl font-semibold">
              Mobile App Authentication (2FA)
            </h3>
            <p className="mb-4">
              Secure your account with TOTP two-factor authentication.
            </p>
            {store.authUser?.otp_enabled ? (
              <button
                type="button"
                className="focus:outline-none text-white bg-purple-700 hover:bg-purple-800 focus:ring-4 focus:ring-purple-300 font-medium rounded-lg text-sm px-5 py-2.5 mb-2"
                onClick={() => disableTwoFactorAuth(user?.id!)}
              >
                Disable 2FA
              </button>
            ) : (
              <button
                type="button"
                className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 mr-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none"
                onClick={() =>
                  generateOTPCode({ user_id: user?.id!})
                }
              >
                Setup 2FA
              </button>
              
            )}
          </div>
        </div>
      </section>
      {openModal && (
        <TwoFactorAuth
          otp_token={secret.otp_token}
          user_id={store.authUser?.id!}
          closeModal={() => setOpenModal(false)}
        />
      )}
    </>
  );
};

export default ProfilePage;
