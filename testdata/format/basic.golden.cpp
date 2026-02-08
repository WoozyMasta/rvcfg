class CfgVehicles
{
  class Health: Health {};
  class Car: Vehicle
  {
    wheels[] = {1, 2, {3, 4}};
    speed = -1;
  };
};
