class Solution:
    def convertTemperature(self, celsius: float) -> List[float]:
        kelvin = celsius + 273.15
        fahrenheit = (9/5*celsius) + 32
        return [kelvin, fahrenheit]