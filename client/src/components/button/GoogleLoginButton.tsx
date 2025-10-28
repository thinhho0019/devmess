// import React from "react";

// interface Props {
//   onLogin: (token: string) => void;
// }

// const GoogleLoginButton: React.FC<Props> = () => {
//   const handleGoogleLogin = () => {
//     const redirectUri = encodeURIComponent(window.location.origin + "/auth/callback");
//     window.location.href = `http://localhost:8080/api/auth/google?redirect_uri=${redirectUri}`;
//   };

//   return (
//     <button
//       onClick={handleGoogleLogin}
//       className="w-full flex items-center justify-center gap-3 bg-[#1f1f1f] text-gray-100 border border-gray-700 
//                  py-2.5 px-4 rounded-xl font-medium shadow-md transition-all duration-200 
//                  hover:bg-[#2a2a2a] hover:shadow-lg active:scale-95"
//     >
//       <img
//         src="https://www.svgrepo.com/show/475656/google-color.svg"
//         alt="Google"
//         className="w-5 h-5"
//       />
//       <span>Đăng nhập với Google</span>
//     </button>
//   );
// };

// export default GoogleLoginButton;
import React from "react";

interface Props {
  onLogin: (token: string) => void;
}

const GoogleLoginButton: React.FC<Props> = ({ onLogin }) => {
  console.log("GoogleLoginButton rendered", onLogin);
  const handleGoogleLogin = () => {
    // Dùng origin của trang hiện tại, sẽ tự động là localhost hoặc devmess.cloud
    // Không cần truyền redirect_uri nếu backend đã dùng env GOOGLE_REDIRECT_URL
    const backendUrl =
      import.meta.env.VITE_API_URL || "http://localhost:8080";

    window.location.href = `${backendUrl}/v1/auth/google`;
  };

  return (
    <button
      onClick={handleGoogleLogin}
      className="w-full flex items-center justify-center gap-3 bg-[#1f1f1f] text-gray-100 border border-gray-700 
                 py-2.5 px-4 rounded-xl font-medium shadow-md transition-all duration-200 
                 hover:bg-[#2a2a2a] hover:shadow-lg active:scale-95"
    >
      <img
        src="https://www.svgrepo.com/show/475656/google-color.svg"
        alt="Google"
        className="w-5 h-5"
      />
      <span>Đăng nhập với Google</span>
    </button>
  );
};

export default GoogleLoginButton;
