function Hero() {
  return (
    <div className="px-8 mt-8 mb-11.75">
      <div className="bg-[url(/images/hero.jpg)] bg-no-repeat bg-cover w-full h-87.5 rounded-xl bg-right flex justify-center items-center">
        <div className="flex justify-center items-center flex-col text-white text-center">
          <h1 className="text-[94px] mb-6 leading-none font-semibold font-hero ">
            Find your favorite <br /> place here!
          </h1>
          <p className="font-[Quicksand] text-[24px] mb-6 leading-none font-medium">
            The best prices for over 2 million properties worldwide
          </p>
        </div>
      </div>
    </div>
  );
}

export default Hero