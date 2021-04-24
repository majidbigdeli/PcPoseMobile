package com.awesomeproject;

import hello.Hello;
import hello.Response;

import com.facebook.react.bridge.ReactApplicationContext;
import com.facebook.react.bridge.ReactContextBaseJavaModule;
import com.facebook.react.bridge.ReactMethod;
import com.google.gson.Gson;

import java.lang.reflect.Type;

public class HelloModule extends ReactContextBaseJavaModule {
    public HelloModule(ReactApplicationContext reactContext) {
        super(reactContext);
    }

    @Override
    public String getName() {
        return "CoreM";
    }

    @ReactMethod(isBlockingSynchronousMethod = true)
    public void Start() {
           Hello.start();
    }

    @ReactMethod(isBlockingSynchronousMethod = true)
    public void Response(String r) {
        Hello.returnResponse(r);
    }

    @ReactMethod(isBlockingSynchronousMethod = true)
    public String GetIp(){
        return  Hello.getOutboundIP();
    }
}
