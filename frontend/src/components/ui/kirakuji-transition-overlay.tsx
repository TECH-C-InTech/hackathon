"use client";

import type { CSSProperties } from "react";
import Image from "next/image";

type KirakujiTransitionOverlayProps = {
  message?: string;
  subMessage?: string;
};

export const KIRAKUJI_TRANSITION_MS = 1600;

/**
 * きらくじ結果へ遷移する間のアニメーション画面を表示する。
 */
export default function KirakujiTransitionOverlay({
  message = "少し待ってて!あなたのためのお告げを探すから。",
  subMessage = "きらくじを引いています...",
}: KirakujiTransitionOverlayProps) {
  const animationStyle = {
    "--kirakuji-transition-duration": `${KIRAKUJI_TRANSITION_MS}ms`,
  } as CSSProperties;

  return (
    <div
      className="fixed inset-0 z-[60] flex items-center justify-center bg-zinc-50 px-4 py-8 font-sans text-zinc-900"
      role="status"
      aria-live="polite"
      aria-busy="true"
      style={animationStyle}
    >
      <div className="relative flex w-full max-w-sm flex-col items-center gap-5 text-center">
        <div className="relative flex items-center justify-center">
          <div
            className="kirakuji-ring absolute h-[220px] w-[220px] rounded-full border border-zinc-200"
            aria-hidden="true"
          />
          <div
            className="kirakuji-ring absolute h-[220px] w-[220px] rounded-full border border-zinc-200"
            style={{ animationDelay: "0.3s" }}
            aria-hidden="true"
          />
          <Image
            src="/kirakuji.png"
            alt=""
            width={240}
            height={240}
            className="relative w-[200px] kirakuji-float md:w-[240px]"
            priority
          />
          <Image
            src="/kirakuji.png"
            alt=""
            width={240}
            height={240}
            className="absolute -z-10 w-[200px] opacity-15 kirakuji-shadow md:w-[240px]"
            aria-hidden="true"
          />
        </div>
        <p className="font-medium text-base">{message}</p>
        <div
          className="h-2 w-full max-w-xs overflow-hidden rounded-full bg-zinc-200"
          role="progressbar"
          aria-label="画面を切り替えています"
        >
          <div className="h-full rounded-full bg-zinc-700 kirakuji-progress" />
        </div>
        <p className="text-sm text-zinc-500">{subMessage}</p>
      </div>
    </div>
  );
}
